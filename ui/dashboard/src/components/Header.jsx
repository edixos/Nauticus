import PropTypes from 'prop-types';

function Header({ title, logoSrc, logoAlt }) {
    return (
        <header className="bg-gray-800 text-white p-4 flex justify-between items-center">
            <div className="flex items-center">
                <img src={logoSrc} alt={logoAlt || "logo"} className="h-10 mr-3" />
                <span className="text-xl font-semibold">{title || ""}</span>
            </div>
            {/* Add any additional header content here */}
        </header>
    );
}

Header.propTypes = {
    logoSrc: PropTypes.string.isRequired,
    title: PropTypes.string,
    logoAlt: PropTypes.string
}

export default Header;