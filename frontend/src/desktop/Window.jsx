import React from 'react';
import Draggable from 'react-draggable';
import { Resizable } from 'react-resizable';

export default function Window({ window, layout, children, onFocus, onResize, onMove, onClose, onMinimize }) {
  const handleMouseDown = () => {
    if (onFocus) {
      onFocus();
    }
    document.body.classList.add('grabbing');
  };

  const handleMouseUp = () => {
    document.body.classList.remove('grabbing');
  };

  const { x, y, width, height, zIndex, minimized } = { ...window, ...layout };

  return (
    <Draggable
      handle=".window-header"
      position={{ x: x || 50, y: y || 50 }}
      onStop={onMove}
      onMouseDown={handleMouseDown}
      onMouseUp={handleMouseUp}
    >
      <Resizable
        width={width || 600}
        height={height || 400}
        onResize={onResize}
        minConstraints={[200, 150]}
        handle={<div className="window-resize-handle" />}
      >
        <div
          className={`window ${minimized ? 'minimized' : ''}`}
          style={{ width: width || 600, height: height || 400, zIndex: zIndex || 'auto' }}
        >
          <div className="window-header">
            <span>{window.title}</span>
            <div>
              <button onClick={onMinimize}>—</button>
              <button onClick={onClose}>×</button>
            </div>
          </div>
          <div className="window-body">{children}</div>
        </div>
      </Resizable>
    </Draggable>
  );
}
