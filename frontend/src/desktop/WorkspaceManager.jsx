import React, { useState } from 'react';
import Window from './Window';

export default function WorkspaceManager({ windows, layout, onLayoutChange, onCloseWindow, onMinimizeWindow, onFocusWindow }) {
  const [activeWorkspace, setActiveWorkspace] = useState(0);

  const workspaces = [0, 1, 2]; // Example: 3 workspaces

  const switchWorkspace = (index) => setActiveWorkspace(index);

  return (
    <>
      {workspaces.map((wsIndex) =>
        wsIndex === activeWorkspace
          ? windows.map((w) => (
              <Window
                key={`${w.id}-ws${wsIndex}`}
                window={w}
                layout={layout[w.id]}
                onMove={(_, { x, y }) => onLayoutChange({ ...layout, [w.id]: { ...layout[w.id], x, y } })}
                onResize={(_, { size }) => onLayoutChange({ ...layout, [w.id]: { ...layout[w.id], width: size.width, height: size.height } })}
                onClose={() => onCloseWindow(w.id)}
                onMinimize={() => onMinimizeWindow(w.id)}
                onFocus={() => onFocusWindow(w.id)}
              >
                {w.component}
              </Window>
            ))
          : null
      )}
      <div className="workspace-switcher">
        {workspaces.map((wsIndex) => (
          <button key={wsIndex} onClick={() => switchWorkspace(wsIndex)}>
            Workspace {wsIndex + 1}
          </button>
        ))}
      </div>
    </>
  );
}
