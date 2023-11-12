import PropTypes from 'prop-types';

function Button({ children, className, icon: Icon, iconClassName, ...props }) {
    const cssClasses = `inline-flex items-center font-bold py-2 px-4 rounded ${className}`;
    const iconClasses = `mr-2  ${iconClassName}`

    return (
        <button className={cssClasses} {...props}>
            {Icon && <Icon className={iconClasses} />} {children}
        </button>
    );
}

Button.propTypes = {
    children: PropTypes.node.isRequired,
    className: PropTypes.string,
    iconClassName: PropTypes.string,
    icon: PropTypes.elementType, // This allows you to pass an icon component type
};

export default Button;