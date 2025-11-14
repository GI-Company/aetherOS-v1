import VFSProxy from './vfsProxy';
import { Server, WebSocket } from 'mock-socket';

describe('VFSProxy', () => {
  let mockServer: Server;
  let vfsProxy: VFSProxy;
  const wsUrl = 'ws://localhost:8080/ws';

  beforeEach(() => {
    mockServer = new Server(wsUrl);
    // Mock the global WebSocket object
    (global as any).WebSocket = WebSocket;
    vfsProxy = new VFSProxy(wsUrl);
  });

  afterEach(() => {
    mockServer.stop();
  });

  test('should send a message to the WebSocket server', (done) => {
    const testTopic = 'test-topic';
    const testPayload = { message: 'hello' };

    mockServer.on('connection', (socket) => {
      socket.on('message', (data) => {
        const message = JSON.parse(data as string);
        expect(message.topic).toBe(testTopic);
        expect(message.payload).toEqual(testPayload);
        done();
      });
    });

    // Allow time for the real VFSProxy to establish a connection
    setTimeout(() => {
      vfsProxy.sendMessage(testTopic, testPayload);
    }, 100);
  });
  
  test('should handle incoming messages from the server', async () => {
    const writePath = 'test/file.txt';
    const writeContent = 'hello world';

    // Let's mock the DB methods for this test to isolate the WebSocket logic
    const mockDb = {
      put: jest.fn().mockResolvedValue(undefined),
      get: jest.fn(),
      getAllKeys: jest.fn(),
    };

    // @ts-ignore - Replace the db instance with our mock
    vfsProxy.db = mockDb;

    const writeMessage = {
      id: '123',
      topic: 'vfs:write',
      payload: { path: writePath, content: writeContent },
    };

    // Trigger the server to send a message to the client
    mockServer.emit('message', JSON.stringify(writeMessage));

    // Give the event loop a moment to process the message
    await new Promise(resolve => setTimeout(resolve, 100));
    
    // Check if the database `put` method was called correctly
    expect(mockDb.put).toHaveBeenCalledWith('files', writeContent, writePath);
  });
});
