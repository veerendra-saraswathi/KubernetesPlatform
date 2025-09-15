import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Sidebar from './components/Sidebar';
import Topbar from './components/Topbar';
import Clusters from './pages/Clusters';
import Workloads from './pages/Workloads';
import Backups from './pages/Backups';
import Settings from './pages/Settings';
import Quotas from './pages/Quotas';

function App() {
  return (
    <Router>
      <div className="flex h-screen bg-gray-100">
        <Sidebar />
        <div className="flex-1 flex flex-col">
          <Topbar />
          <main className="flex-1 p-4 overflow-auto">
            <Routes>
              <Route path="/" element={<Clusters />} />
              <Route path="/workloads" element={<Workloads />} />
              <Route path="/backups" element={<Backups />} />
              <Route path="/settings" element={<Settings />} />
              <Route path="/quotas" element={<Quotas />} />
            </Routes>
          </main>
        </div>
      </div>
    </Router>
  );
}

export default App;

