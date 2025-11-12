import React from 'react';
import Draggable from 'react-draggable';
import { Resizable } from 'react-resizable';

export default function Window({ window, onFocus, onResize, onMove, onClose }) {
  const handleMouseDown = () => {
    onFocus();
    document.body.classList.add('grabbing');
  };

  const handleMouseUp = () => {
    document.body.classList.remove('grabbing');
  };

  return (
    <Draggable
      handle=".window-header"
      position={{ x: window.x, y: window.y }}
      onStop={onMove}
      onMouseDown={handleMouseDown}
      onMouseUp={handleMouseUp}
    >
      <Resizable
        width={window.width}
        height={window.height}
        onResize={onResize}
        minConstraints={[200, 150]}
        handle={<div className="window-resize-handle" />}
      >
        <div
          className={`window ${window.minimized ? 'minimized' : ''}`}
          style={{ width: window.width, height: window.height, zIndex: window.zIndex }}
        >
          <div className="window-header">
            <span>{window.title}</span>
            <button onClick={onClose}>×</button>
          </div>
          <div className="window-body">{window.children}</div>
        </div>
      </Resizable>
    </Draggable>
  );
}
