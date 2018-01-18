package p2p

import (
	"github.com/gogo/protobuf/proto"
	"github.com/spacemeshos/go-spacemesh/crypto"
	"github.com/spacemeshos/go-spacemesh/log"
	"github.com/spacemeshos/go-spacemesh/p2p/dht"
	"github.com/spacemeshos/go-spacemesh/p2p/node"
	"github.com/spacemeshos/go-spacemesh/p2p/nodeconfig"
	"github.com/spacemeshos/go-spacemesh/p2p/pb"
)

// The local Spacemesh node is the root of all evil
type LocalNode interface {
	Id() []byte
	String() string
	Pretty() string

	PrivateKey() crypto.PrivateKey
	PublicKey() crypto.PublicKey

	DhtId() dht.ID
	TcpAddress() string

	Sign(data proto.Message) ([]byte, error)
	SignToString(data proto.Message) (string, error)
	NewProtocolMessageMetadata(protocol string, reqId []byte, gossip bool) *pb.Metadata

	GetSwarm() Swarm
	GetPing() Ping

	Config() nodeconfig.Config

	GetRemoteNodeData() node.RemoteNodeData

	Shutdown()

	// logging wrappers - log node id and args

	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Error(format string, args ...interface{})
	Warning(format string, args ...interface{})

	// local store persistence
	ensureNodeDataDirectory() (string, error)
	persistData() error
}

// Creates a local node with a provided tcp address
// Attempts to set node identity from persisted data in local store
// Creates a new identity if none was loads
func NewLocalNode(tcpAddress string, config nodeconfig.Config, persist bool) (LocalNode, error) {

	if len(nodeconfig.ConfigValues.NodeId) > 0 {
		// user provided node id/pubkey via the cli - attempt to start that node w persisted data
		data, err := readNodeData(nodeconfig.ConfigValues.NodeId)
		if err != nil {
			return nil, err
		}

		return newNodeFromData(tcpAddress, data, config, persist)
	}

	// look for persisted node data in the nodes directory
	// load the node with the data of the first node found
	nodeData, err := readFirstNodeData()
	if err != nil {
		return nil, err
	}

	if nodeData != nil {
		// crete node using persisted node data
		return newNodeFromData(tcpAddress, nodeData, config, persist)
	}

	// generate new node
	return NewNodeIdentity(tcpAddress, config, persist)
}

// Creates a new local node without attempting to restore identity from local store
func NewNodeIdentity(tcpAddress string, config nodeconfig.Config, persist bool) (LocalNode, error) {
	priv, pub, _ := crypto.GenerateKeyPair()
	return newLocalNodeWithKeys(pub, priv, tcpAddress, config, persist)
}

func newLocalNodeWithKeys(pubKey crypto.PublicKey, privKey crypto.PrivateKey, tcpAddress string,
	config nodeconfig.Config, persist bool) (LocalNode, error) {

	n := &localNodeImp{
		pubKey:     pubKey,
		privKey:    privKey,
		tcpAddress: tcpAddress,
		config:     config, // store this node passed-in config values and use them later
		dhtId:      dht.NewIdFromNodeKey(pubKey.Bytes()),
	}

	dataDir, err := n.ensureNodeDataDirectory()
	if err != nil {
		return nil, err
	}

	// setup logging
	n.logger = log.CreateLogger(n.pubKey.Pretty(), dataDir, "node.log")

	// swarm owned by node
	s, err := NewSwarm(tcpAddress, n)
	if err != nil {
		n.Error("can't create a local node without a swarm", err)
		return nil, err
	}

	n.swarm = s
	n.ping = NewPingProtocol(s)

	if persist {
		// persist store data so we can start it on future app sessions
		err = n.persistData()
		if err != nil { // no much use of starting if we can't store node private key in store
			n.Error("failed to persist node data to local store", err)
			return nil, err
		}
	}

	return n, nil
}

// Creates a new node from peristed NodeData
func newNodeFromData(tcpAddress string, d *NodeData, config nodeconfig.Config, persist bool) (LocalNode, error) {
	priv := crypto.NewPrivateKeyFromString(d.PrivKey)
	pub, err := crypto.NewPublicKeyFromString(d.PubKey)
	if err != nil {
		log.Error("failed to create public key from string", err)
		return nil, err
	}

	return newLocalNodeWithKeys(pub, priv, tcpAddress, config, persist)
}
