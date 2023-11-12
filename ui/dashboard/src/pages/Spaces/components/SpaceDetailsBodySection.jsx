import { PropTypes } from 'prop-types';


function SpaceDetailsBodySection({ title, children }) {
    return (
        < div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6" >
            <dt className="text-sm font-medium text-gray-500">
                {title}
            </dt>
            <dd className="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                {children}
            </dd>
        </div >
    );
}

SpaceDetailsBodySection.propTypes = {
    children: PropTypes.node.isRequired,
    title: PropTypes.string.isRequired,
}


export default SpaceDetailsBodySection;