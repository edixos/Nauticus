import { PropTypes } from 'prop-types';

function Label({ children, className, ...props }) {
    const cssClasses = `${className} inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium  m-1`;

    return (
        <span {...props} className={cssClasses}>
            {children}
        </span>
    );
}

Label.propTypes = {
    children: PropTypes.node.isRequired,
    className: PropTypes.string,
}

export default Label;