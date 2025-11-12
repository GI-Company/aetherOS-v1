import React, { useState } from 'react';
import Window from './Window';

export default function WorkspaceManager({ windows, layout, onLayoutChange }) {
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
                title={w.title}
                layout={layout[w.id]}
                onLayoutChange={(newLayout) => onLayoutChange({ ...layout, [w.id]: newLayout })}
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