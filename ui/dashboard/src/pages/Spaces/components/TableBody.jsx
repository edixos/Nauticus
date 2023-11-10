import { PropTypes } from 'prop-types';

function TableBody({ data }) {
    return (
        <tbody className="text-gray-100">
            {data.map((space) => (
                <tr key={space.metadata.name}>
                    <td className="px-6 py-4 border-b border-gray-200">
                        {space.metadata.name}
                    </td>
                    <td className="px-6 py-4 border-b border-gray-200">
                        {space.spec.owners.reduce((fistOwner, owner) => fistOwner.kind + ':' + fistOwner.name + '::' + owner.kind + ':' + owner.name)}
                    </td>
                    <td className="px-6 py-4 border-b border-gray-200">
                        {space.spec.resourceQuota.hard["limits.cpu"]}
                    </td>
                    <td className="px-6 py-4 border-b border-gray-200">
                        {space.spec.resourceQuota.hard["limits.memory"]}
                    </td>

                </tr>
            ))}
        </tbody>
    );
}

TableBody.propTypes = {
    data: PropTypes.array.isRequired,
}

export default TableBody;