
import React, { useEffect, useState } from 'react';
import Desktop from './desktop/Desktop';
import { AIAgentWindow, ComputeWindow, NotificationWindow, MarketplaceWindow } from './windows';
import { restoreLayout, saveLayout } from './sdk/layoutManager';

export default function App() {
  const [layout, setLayout] = useState({});
  const [windows, setWindows] = useState([]);

  useEffect(() => {
    setWindows([
      { id: 'ai', title: 'AI Agent', component: <AIAgentWindow /> },
      { id: 'compute', title: 'Compute', component: <ComputeWindow /> },
      { id: 'notifications', title: 'Notifications', component: <NotificationWindow /> },
      { id: 'marketplace', title: 'Marketplace', component: <MarketplaceWindow /> }
    ]);

    restoreLayout((savedLayout) => setLayout(savedLayout));
  }, []);

  const handleLayoutChange = (newLayout) => {
    setLayout(newLayout);
    saveLayout(newLayout);
  };

  return <Desktop windows={windows} layout={layout} onLayoutChange={handleLayoutChange} />;

}
