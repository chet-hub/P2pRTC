using System;
using System.Runtime.InteropServices;


public class P2PSocket
    {
        const string dll = "\pathto\webrtc.so";

        public delegate void GoOnSignalCb(UInt64 p2pSocket, string message);

        public delegate void GoOnOpenCb(UInt64 p2pSocket);

        public delegate void GoOnCloseCb(UInt64 p2pSocket);

        public unsafe delegate  void GoOnMessageCb(UInt64 p2pSocket, byte* message, UInt64 length);
        
        public delegate void OnMessageCb(UInt64 p2pSocket, byte[] message);

        public delegate void GoOnErrorCb(UInt64 p2pSocket, string message);

        [DllImport(dll)]
        private static extern UInt64 newClient(int cIsTrickleIce, string cWebrtcconfig, string cDatachannelconfig);

        [DllImport(dll)]
        private static extern UInt64 newServer(int cIsTrickleIce, string cWebrtcconfig, string cDatachannelconfig);

        [DllImport(dll)]
        private static extern int signal(UInt64 p2pSocket, string message);

        [DllImport(dll)]
        private static extern int connect(UInt64 p2pSocket);

        [DllImport(dll)]
        private static extern int closed(UInt64 p2pSocket);

        [DllImport(dll)]
        private static extern int opened(UInt64 p2pSocket);

        [DllImport(dll)]
        private static extern int connecting(UInt64 p2pSocket);

        [DllImport(dll)]
        private static extern int ordered(UInt64 p2pSocket);

        [DllImport(dll)]
        private static extern int close(UInt64 p2pSocket);

        [DllImport(dll)]
        private static extern int send(UInt64 p2pSocket, byte[] p, int length);

        [DllImport(dll)]
        private static extern int listenOnSignal(UInt64 p2pSocket, GoOnSignalCb cb);

        [DllImport(dll)]
        private static extern int listenOnError(UInt64 p2pSocket, GoOnErrorCb cb);

        [DllImport(dll)]
        private static extern int listenOnMessage(UInt64 p2pSocket, GoOnMessageCb cb);

        [DllImport(dll)]
        private static extern int listenOnOpen(UInt64 p2pSocket, GoOnOpenCb cb);

        [DllImport(dll)]
        private static extern int listenOnClose(UInt64 p2pSocket, GoOnCloseCb cb);


        private UInt64 connection;

        private P2PSocket(bool isServer, bool isTrickleIce, string webRtcConfig, string dataChannelConfig)
        {
            if (isServer)
            {
                connection = newServer(isTrickleIce ? 1 : 0, webRtcConfig, dataChannelConfig);
            }
            else
            {
                connection = newClient(isTrickleIce ? 1 : 0, webRtcConfig, dataChannelConfig);
            }
        }

        public static P2PSocket NewClient(bool isTrickleIce, string webRtcConfig, string dataChannelConfig)
        {
            return new P2PSocket(false, isTrickleIce, webRtcConfig, dataChannelConfig);
        }

        public static P2PSocket NewServer(bool isTrickleIce, string webRtcConfig, string dataChannelConfig)
        {
            return new P2PSocket(true, isTrickleIce, webRtcConfig, dataChannelConfig);
        }

        public void ListenOnSignal(GoOnSignalCb cb)
        {
            if (listenOnSignal(this.connection, cb) == 0)
            {
                throw new Exception("Can't find the connection");
            }
        }

        public void ListenOnError(GoOnErrorCb cb)
        {
            if (listenOnError(this.connection, cb) == 0)
            {
                throw new Exception("Can't find the connection");
            }
        }

        public void ListenOnMessage(OnMessageCb cb)
        {
            unsafe
            {
                GoOnMessageCb goOnMessageCb = (UInt64 p2pSocket, byte* message, UInt64 length) =>
                {
                    //todo 
                    byte[] data = new byte[length];
                    for (ulong i = 0; i < length; i++)
                    {
                        data[i] = *(message+i);
                    }
                    cb(p2pSocket, data);
                };
                if (listenOnMessage(this.connection, goOnMessageCb) == 0)
                {
                    throw new Exception("Can't find the connection");
                }
            }
        }

        public void ListenOnOpen(GoOnOpenCb cb)
        {
            if (listenOnOpen(this.connection, cb) == 0)
            {
                throw new Exception("Can't find the connection");
            }
        }

        public void ListenOnClose(GoOnCloseCb cb)
        {
            if (listenOnClose(this.connection, cb) == 0)
            {
                throw new Exception("Can't find the connection");
            }
        }

        public bool Signal(string message)
        {
            int result = signal(this.connection, message);
            if (result == -1)
            {
                throw new Exception("Can't find the connection");
            }
            else
            {
                return result == 1;
            }
        }

        public bool Connect()
        {
            int result = connect(this.connection);
            if (result == -1)
            {
                throw new Exception("Can't find the connection");
            }
            else
            {
                return result == 1;
            }
        }

        public bool Closed()
        {
            int result = closed(this.connection);
            if (result == -1)
            {
                throw new Exception("Can't find the connection");
            }
            else
            {
                return result == 1;
            }
        }

        public bool Opened()
        {
            int result = opened(this.connection);
            if (result == -1)
            {
                throw new Exception("Can't find the connection");
            }
            else
            {
                return result == 1;
            }
        }

        public bool IsConnecting()
        {
            int result = connecting(this.connection);
            if (result == -1)
            {
                throw new Exception("Can't find the connection");
            }
            else
            {
                return result == 1;
            }
        }

        public bool IsOrdered()
        {
            int result = ordered(this.connection);
            if (result == -1)
            {
                throw new Exception("Can't find the connection");
            }
            else
            {
                return result == 1;
            }
        }

        public bool Close()
        {
            int result = close(this.connection);
            if (result == -1)
            {
                throw new Exception("Can't find the connection");
            }
            else
            {
                return result == 1;
            }
        }


        public bool Send(byte[] msg)
        {
            int result = send(this.connection, msg, msg.Length);
            if (result == -1)
            {
                throw new Exception("Can't find the connection");
            }
            else
            {
                return result == 1;
            }
        }
    }
