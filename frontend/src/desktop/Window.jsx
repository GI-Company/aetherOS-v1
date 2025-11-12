import React, { useState, useRef, useEffect } from 'react';
import '../styles/desktop.css';

export default function Window({ title, children, layout, onLayoutChange }) {
  const windowRef = useRef(null);
  const [position, setPosition] = useState(layout?.position || { x: 100, y: 100 });
  const [size, setSize] = useState(layout?.size || { width: 400, height: 300 });
  const [minimized, setMinimized] = useState(false);
  const [dragging, setDragging] = useState(false);
  const [resizing, setResizing] = useState(false);
  const [offset, setOffset] = useState({ x: 0, y: 0 });

  const toggleMinimize = () => setMinimized(!minimized);

  // Dragging
  const onMouseDownHeader = (e) => {
    setDragging(true);
    setOffset({ x: e.clientX - position.x, y: e.clientY - position.y });
  };

  const onMouseMove = (e) => {
    if (dragging) {
      let newX = e.clientX - offset.x;
      let newY = e.clientY - offset.y;

      // Snap to edges
      const snapMargin = 20;
      if (newX < snapMargin) newX = 0;
      if (newY < snapMargin) newY = 0;
      if (window.innerWidth - (newX + size.width) < snapMargin) newX = window.innerWidth - size.width;
      if (window.innerHeight - (newY + size.height + 48) < snapMargin) newY = window.innerHeight - size.height - 48; // 48px taskbar

      setPosition({ x: newX, y: newY });
      onLayoutChange({ position: { x: newX, y: newY }, size });
    }

    if (resizing) {
      const newWidth = Math.max(300, e.clientX - position.x);
      const newHeight = Math.max(200, e.clientY - position.y);
      setSize({ width: newWidth, height: newHeight });
      onLayoutChange({ position, size: { width: newWidth, height: newHeight } });
    }
  };

  const onMouseUp = () => {
    setDragging(false);
    setResizing(false);
  };

  useEffect(() => {
    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
    return () => {
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };
  });

  return (
    <div
      ref={windowRef}
      className={`window glassmorphic ${minimized ? 'minimized' : ''}`}
      style={{ top: position.y, left: position.x, width: size.width, height: size.height }}
    >
      <div className="window-header" onMouseDown={onMouseDownHeader}>
        <span>{title}</span>
        <button onClick={toggleMinimize}>{minimized ? '🔼' : '🔽'}</button>
      </div>
      {!minimized && (
        <div className="window-body">
          {children}
          <div
            className="window-resize-handle"
            onMouseDown={() => setResizing(true)}
          />
        </div>
      )}
    </div>
  );
}