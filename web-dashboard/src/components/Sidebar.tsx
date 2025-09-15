import React from 'react';
import { Link } from 'react-router-dom';

const Sidebar = () => {
  return (
    <aside className="w-64 bg-gray-800 text-white flex flex-col p-4">
      <h2 className="text-2xl font-bold mb-4">GPU Operator</h2>
      <nav className="flex flex-col space-y-2">
        <Link to="/" className="hover:text-gray-300">Clusters</Link>
        <Link to="/workloads" className="hover:text-gray-300">Workloads</Link>
        <Link to="/backups" className="hover:text-gray-300">Backups</Link>
        <Link to="/settings" className="hover:text-gray-300">Settings</Link>
      </nav>
    </aside>
  );
};

export default Sidebar;

