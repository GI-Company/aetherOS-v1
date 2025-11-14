import { openDB, IDBPDatabase } from 'idb';

class VFSProxy {
  private ws: WebSocket;
  private db: IDBPDatabase | null = null;

  constructor(private wsUrl: string) {
    this.ws = new WebSocket(wsUrl);
    this.initDB();
    this.initWebSocket();
  }

  private async initDB() {
    this.db = await openDB('AetherVFS', 1, {
      upgrade(db) {
        db.createObjectStore('files');
      },
    });
  }

  private initWebSocket() {
    this.ws.onopen = () => {
      console.log('VFSProxy connected to kernel');
    };

    this.ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      this.handleKernelMessage(msg);
    };

    this.ws.onclose = () => {
      console.log('VFSProxy disconnected from kernel');
      // Implement reconnection logic here
    };

    this.ws.onerror = (error) => {
      console.error('VFSProxy WebSocket error:', error);
    };
  }

  private async handleKernelMessage(msg: any) {
    if (!this.db) return;

    const { topic, payload } = msg;

    switch (topic) {
      case 'vfs:write':
        await this.db.put('files', payload.content, payload.path);
        this.sendAck(msg.id, 'success');
        break;
      case 'vfs:read':
        const content = await this.db.get('files', payload.path);
        this.sendMessage('vfs:read:resp', { path: payload.path, content });
        break;
      case 'vfs:list':
        const keys = await this.db.getAllKeys('files');
        this.sendMessage('vfs:list:resp', { keys });
        break;
    }
  }

  private sendMessage(topic: string, payload: any) {
    this.ws.send(JSON.stringify({ topic, payload }));
  }

  private sendAck(id: string, status: string) {
    this.sendMessage('ack', { id, status });
  }
}

export default VFSProxy;
