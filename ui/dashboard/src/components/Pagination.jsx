import PropTypes, { number } from 'prop-types';
import { useEffect, useState } from 'react';

function Pagination({
    pageCount,
    pageIndex,
    initialPageSize,
    canPreviousPage,
    canNextPage,
    gotoPage,
    nextPage,
    previousPage,
    setPageSize: changePageSize,
    maxPageGroup = 6,
    pageSizeOptions = [1, 20, 30, 40, 50]
}) {
    // State for the current page size
    const [pageSize, setPageSize] = useState(initialPageSize);

    // Effect to update the page size in the parent component when it changes here
    useEffect(() => {
        changePageSize(pageSize);
    }, [pageSize, changePageSize]);

    // Maximum number of pages in a group
    const currentPageGroup = Math.floor(pageIndex / maxPageGroup);
    const startPage = currentPageGroup * maxPageGroup;
    const endPage = Math.min(startPage + maxPageGroup, pageCount);


    // Generate page number buttons
    // Calculate the range of pages in the current group
    const pageNumbers = Array.from(
        { length: endPage - startPage },
        (_, index) => startPage + index
    );

    return (
        <div className="flex items-center space-x-1">
            {/* Page size selection */}
            <select
                value={pageSize}
                onChange={e => setPageSize(Number(e.target.value))}
                className="border-2 border-gray-300 text-gray-700 rounded-lg p-1"
            >
                {pageSizeOptions.map(size => (
                    <option key={size} value={size}>
                        {size} per page
                    </option>
                ))}
            </select>

            {/* Pagination controls */}
            <button
                onClick={() => gotoPage(0)}
                disabled={!canPreviousPage}
                className="px-3 py-1 rounded-full bg-gray-200 text-gray-500 hover:bg-blue-500 hover:text-white disabled:opacity-10"
            >
                {'<<'}
            </button>
            <button
                onClick={() => previousPage()}
                disabled={!canPreviousPage}
                className="px-3 py-1 rounded-full bg-gray-200 text-gray-500 hover:bg-blue-500 hover:text-white disabled:opacity-10"
            >
                {'<'}
            </button>

            {/* Show "..." button to go to previous group if needed */}
            {currentPageGroup > 0 && (
                <button
                    onClick={() => gotoPage((currentPageGroup - 1) * maxPageGroup)}
                    className="px-3 py-1 rounded-full bg-gray-200 text-gray-500 hover:bg-blue-500 hover:text-white"
                >
                    {'...'}
                </button>
            )}

            {/* Page number buttons */}
            {pageNumbers.map(number => (
                <button
                    key={number}
                    onClick={() => gotoPage(number)}
                    className={`px-3 py-1 rounded-full ${pageIndex === number ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-500 hover:bg-blue-500 hover:text-white'}`}
                >
                    {number + 1}
                </button>
            ))}

            {/* Show "..." button to go to next group if needed */}
            {currentPageGroup < Math.ceil(pageCount / maxPageGroup) - 1 && (
                <button
                    onClick={() => gotoPage((currentPageGroup + 1) * maxPageGroup)}
                    className="px-3 py-1 rounded-full bg-gray-200 text-gray-500 hover:bg-blue-500 hover:text-white"
                >
                    {'...'}
                </button>
            )}

            <button
                onClick={() => nextPage()}
                disabled={!canNextPage}
                className="px-3 py-1 rounded-full bg-gray-200 text-gray-500 hover:bg-blue-500 hover:text-white disabled:opacity-10"
            >
                {'>'}
            </button>
            <button
                onClick={() => gotoPage(pageCount - 1)}
                disabled={!canNextPage}
                className="px-3 py-1 rounded-full bg-gray-200 text-gray-500 hover:bg-blue-500 hover:text-white disabled:opacity-10"
            >
                {'>>'}
            </button>
        </div>
    );
}


Pagination.propTypes = {
    pageCount: PropTypes.number.isRequired,
    pageIndex: PropTypes.number.isRequired,
    initialPageSize: PropTypes.number.isRequired,
    canPreviousPage: PropTypes.bool.isRequired,
    canNextPage: PropTypes.bool.isRequired,
    gotoPage: PropTypes.func.isRequired,
    nextPage: PropTypes.func.isRequired,
    previousPage: PropTypes.func.isRequired,
    setPageSize: PropTypes.func.isRequired,
    maxPageGroup: PropTypes.number,
    pageSizeOptions: PropTypes.arrayOf(number)
};
export default Pagination;