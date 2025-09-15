import React from 'react';

interface CatalogItemProps {
  template: {
    id: number;
    name: string;
    description: string;
  };
}

const CatalogItem: React.FC<CatalogItemProps> = ({ template }) => {
  return (
    <div className="border rounded-lg p-4 shadow hover:shadow-lg transition">
      <h2 className="text-xl font-semibold mb-2">{template.name}</h2>
      <p className="text-gray-600 mb-4">{template.description}</p>
      <button className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
        Deploy
      </button>
    </div>
  );
};

export default CatalogItem;

