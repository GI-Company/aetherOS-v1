const socket = new WebSocket("ws://localhost:8080/ws");

const listeners = new Map();

socket.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (listeners.has(message.type)) {
    listeners.get(message.type).forEach((callback) => callback(message.payload));
  }
};

export const send = (type, payload) => {
  socket.send(JSON.stringify({ type, payload }));
};

export const on = (type, callback) => {
  if (!listeners.has(type)) {
    listeners.set(type, []);
  }
  listeners.get(type).push(callback);
};

export const off = (type, callback) => {
  if (listeners.has(type)) {
    const newListeners = listeners.get(type).filter((cb) => cb !== callback);
    listeners.set(type, newListeners);
  }
};

export default socket;
