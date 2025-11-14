import VFSProxy from './vfsProxy';

describe('VFS E2E Test', () => {
  let vfsProxy: VFSProxy;

  beforeAll(() => {
    // Note: This test requires a running Aether kernel
    vfsProxy = new VFSProxy('ws://localhost:8080/ws');
  });

  test('should write and read a file', async () => {
    const filePath = 'test.txt';
    const fileContent = 'Hello, Aether!';

    // Give the WebSocket time to connect
    await new Promise(resolve => setTimeout(resolve, 1000));

    // @ts-ignore - sendMessage is private, but we need it for testing
    vfsProxy.sendMessage('vfs:write', { path: filePath, content: fileContent });

    // Listen for the write acknowledgement
    const ackPromise = new Promise(resolve => {
      // @ts-ignore - ws is private
      vfsProxy.ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        if (msg.topic === 'ack') {
          resolve(msg);
        }
      };
    });

    await expect(ackPromise).resolves.toEqual(expect.objectContaining({ payload: { status: 'success' } }));

    // Now, read the file back
    // @ts-ignore - sendMessage is private
    vfsProxy.sendMessage('vfs:read', { path: filePath });

    const readPromise = new Promise(resolve => {
      // @ts-ignore - ws is private
      vfsProxy.ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        if (msg.topic === 'vfs:read:resp') {
          resolve(msg);
        }
      };
    });

    await expect(readPromise).resolves.toEqual(expect.objectContaining({ payload: { path: filePath, content: fileContent } }));
  });
});
