import React from "react";

const Quotas = () => {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Resource Quotas</h1>
      <table className="table-auto w-full border">
        <thead>
          <tr>
            <th className="border px-4 py-2">Namespace</th>
            <th className="border px-4 py-2">CPU</th>
            <th className="border px-4 py-2">Memory</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td className="border px-4 py-2">tenant1</td>
            <td className="border px-4 py-2">4</td>
            <td className="border px-4 py-2">16Gi</td>
          </tr>
          <tr>
            <td className="border px-4 py-2">tenant2</td>
            <td className="border px-4 py-2">2</td>
            <td className="border px-4 py-2">8Gi</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

export default Quotas;

