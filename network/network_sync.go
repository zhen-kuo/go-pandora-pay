package network

import (
	"pandora-pay/config"
	"pandora-pay/gui"
	"pandora-pay/network/api/api_websockets/consensus"
	"pandora-pay/network/websocks"
	"pandora-pay/network/websocks/connection"
	"pandora-pay/recovery"
	"time"
)

func (network *Network) continuouslyConnectNewPeers() {

	recovery.SafeGo(func() {

		for {

			if network.Websockets.GetClients() >= config.WEBSOCKETS_NETWORK_CLIENTS_MAX {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			knownNode := network.KnownNodes.GetRandomKnownNode()
			if knownNode == nil {
				continue
			}

			if network.BannedNodes.IsBanned(knownNode.URL) {
				continue //banned already
			}

			_, exists := network.Websockets.AllAddresses.Load(knownNode.URL)
			if !exists {

				if knownNode != nil {
					_, err := websocks.NewWebsocketClient(network.Websockets, knownNode)
					if err != nil {

						if err.Error() != "Already connected" {
							if knownNode.DecreaseScore(-5, false) {
								network.KnownNodes.RemoveKnownNode(knownNode)
							}

						}

					} else {
						gui.GUI.Log("connected to: " + knownNode.URL)
					}
				}
			}

			time.Sleep(100 * time.Millisecond)
		}

	})

}

func (network *Network) continuouslyDownloadChain() {
	recovery.SafeGo(func() {

		for {

			if conn := network.Websockets.GetRandomSocket(); conn != nil {
				data, err := connection.SendJSONAwaitAnswer[consensus.ChainUpdateNotification](conn, []byte("get-chain"), nil, nil, 0)
				if err == nil {
					network.Websockets.ApiWebsockets.Consensus.ChainUpdateProcess(conn, data)
				}
			}

			time.Sleep(2000 * time.Millisecond)
		}

	})
}

func (network *Network) continuouslyDownloadMempool() {

	recovery.SafeGo(func() {

		for {

			if conn := network.Websockets.GetRandomSocket(); conn != nil {
				if config.CONSENSUS == config.CONSENSUS_TYPE_FULL && conn.Handshake.Consensus == config.CONSENSUS_TYPE_FULL {
					network.MempoolSync.DownloadMempool(conn)
				}
			}

			time.Sleep(2000 * time.Millisecond)
		}

	})

}

func (network *Network) continuouslyDownloadNetworkNodes() {

	recovery.SafeGo(func() {

		for {

			conn := network.Websockets.GetRandomSocket()
			if conn != nil {

				if config.CONSENSUS == config.CONSENSUS_TYPE_FULL && conn.Handshake.Consensus == config.CONSENSUS_TYPE_FULL {
					network.KnownNodesSync.DownloadNetworkNodes(conn)
				}

			}

			time.Sleep(10000 * time.Millisecond)
		}

	})

}

func (network *Network) syncBlockchainNewConnections() {
	recovery.SafeGo(func() {

		cn := network.Websockets.UpdateNewConnectionMulticast.AddListener()
		defer network.Websockets.UpdateNewConnectionMulticast.RemoveChannel(cn)

		for {

			conn, ok := <-cn
			if !ok {
				return
			}

			//making it async
			recovery.SafeGo(func() {

				data, err := connection.SendJSONAwaitAnswer[consensus.ChainUpdateNotification](conn, []byte("get-chain"), nil, nil, 0)
				if err == nil {
					network.Websockets.ApiWebsockets.Consensus.ChainUpdateProcess(conn, data)
				}

			})

		}
	})
}
