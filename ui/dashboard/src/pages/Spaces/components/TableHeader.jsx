import { PropTypes } from 'prop-types'
import TableHeaderItem from "./TableHeaderItem";

function TableHeader({ headers }) {
    return (
        <thead className="bg-gray-50">
            <tr>
                {headers.map((header) => (<TableHeaderItem key={header.name} name={header.name} />))}
            </tr>
        </thead>
    );
}

TableHeader.propTypes = {
    headers: PropTypes.array.isRequired,
}

export default TableHeader;