import React from 'react';
import CatalogItem from '../../components/catalog/CatalogItem';

const catalogTemplates = [
  { id: 1, name: 'GPU Job - Training', description: 'Submit a training job on GPU cluster.' },
  { id: 2, name: 'Backup Workflow', description: 'Trigger cluster backup.' },
  { id: 3, name: 'Cluster Scale Template', description: 'Add/remove nodes automatically.' },
];

const Catalog: React.FC = () => {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Template Marketplace</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {catalogTemplates.map(template => (
          <CatalogItem key={template.id} template={template} />
        ))}
      </div>
    </div>
  );
};

export default Catalog;

