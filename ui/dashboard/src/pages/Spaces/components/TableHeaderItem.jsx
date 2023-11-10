import { PropTypes } from 'prop-types';

function TableHeaderItem({ name }) {
    return (
        <th className="px-6 py-3 border-b-2 border-gray-200 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
            {name}
        </th>
    );
}

TableHeaderItem.propTypes = {
    name: PropTypes.string,
}

export default TableHeaderItem;