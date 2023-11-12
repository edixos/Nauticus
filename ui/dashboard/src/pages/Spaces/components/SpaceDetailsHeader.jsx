import Button from '../../../components/Button';
import { PropTypes } from 'prop-types';
import { FaRegClone } from 'react-icons/fa6';
import { FaEdit, FaTrash } from 'react-icons/fa';


function SpaceDetailsHeader({ spaceDetails }) {
    return (
        <div className="px-4 py-5 sm:px-6 flex justify-between items-center">
            <h2 className="text-3xl leading-6 font-medium text-gray-900">
                Space: <strong>{spaceDetails.metadata.name}</strong>
            </h2>
            <div className="flex gap-3">
                {/* Action Buttons */}
                <Button icon={FaRegClone} className="bg-blue-500 hover:bg-blue-700 text-white">Clone</Button>
                <Button icon={FaEdit} className="bg-yellow-500 hover:bg-yellow-700 text-white">Edit YAML</Button>
                <Button icon={FaTrash} className="bg-red-500 hover:bg-red-700 text-white">Delete</Button>
            </div>
        </div>
    );
}

SpaceDetailsHeader.propTypes = {
    spaceDetails: PropTypes.object,
}



export default SpaceDetailsHeader;