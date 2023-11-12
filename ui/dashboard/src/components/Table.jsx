import {
    useReactTable, flexRender, getCoreRowModel, getFilteredRowModel, getPaginationRowModel
} from '@tanstack/react-table';
import { useState } from 'react';
import { PropTypes } from 'prop-types';
import { FiSearch } from 'react-icons/fi';
import { BsPlusSquare } from 'react-icons/bs';
import Pagination from './Pagination';
import { Link } from 'react-router-dom';

function Table({ columns, data }) {
    const [pageIndex] = useState(0);
    const [pageSize] = useState(2);
    const [filter, setFilter] = useState('');

    const table = useReactTable({
        columns,
        data,
        state: {
            globalFilter: filter,
            pageIndex, // Current page index
            pageSize
        },
        initialState: {
            pagination: {
                pageSize,
                pageIndex
            }
        },
        onGlobalFilterChange: setFilter,
        globalFilterFn: 'includesString',
        getCoreRowModel: getCoreRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        getPaginationRowModel: getPaginationRowModel()
    });

    console.log(table.getState());

    return (
        <div className='flex flex-col gap-4 p-4'>
            <div className="flex w-full gap-4">
                <div className="flex-grow relative w-2/3">
                    <FiSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-lg text-gray-500" />
                    <input
                        type="text"
                        value={filter}
                        onChange={e => setFilter(e.target.value)}
                        placeholder='Search Spaces ...'
                        className="pl-10 pr-4 py-2 w-full border-2 border-gray-300 rounded-lg text-gray-700 focus:ring-2 focus:ring-blue-500"
                    />
                </div>
                <button className="flex-grow-0 flex items-center gap-2 bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg transition-colors">
                    <BsPlusSquare className="md:text-2xl sm:text-sm" />
                    Add Space
                </button>
            </div>
            <div className="overflow-x-auto rounded-lg shadow">
                <table className="min-w-full divide-y divide-gray-700">
                    <thead className="bg-gray-700">
                        {table.getHeaderGroups().map(headerGroup => (
                            <tr key={headerGroup.id}>
                                {headerGroup.headers.map(header => (
                                    <th key={header.id} scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        {header.isPlaceholder
                                            ? null
                                            : flexRender(header.column.columnDef.header, header.getContext())}
                                    </th>
                                ))}
                            </tr>
                        ))}
                    </thead>
                    <tbody className="bg-gray-200 divide-y divide-gray-900">
                        {table.getPaginationRowModel().rows.map(row => ( // Use the paginated rows here
                            <tr key={row.id}>
                                {row.getVisibleCells().map(cell => {
                                    if (cell.column.id === 'Name') {

                                        return (
                                            <td key={cell.id} className="px-6 py-4 whitespace-nowrap">
                                                <Link to={`${row.original.metadata.name}`} className="text-blue-600 hover:text-blue-800">
                                                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                                </Link>
                                            </td>
                                        );
                                    } else {
                                        // For all other cells, just render the cell value
                                        return (
                                            <td key={cell.id} className="px-6 py-4 whitespace-nowrap">
                                                {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                            </td>
                                        );
                                    }
                                }
                                )}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            <Pagination
                pageCount={table.getPageCount()}
                pageIndex={table.getState().pagination.pageIndex}
                initialPageSize={3}
                canPreviousPage={table.getCanPreviousPage()}
                canNextPage={table.getCanNextPage()}
                gotoPage={table.setPageIndex}
                nextPage={table.nextPage}
                previousPage={table.previousPage}
                setPageSize={table.setPageSize}
            />
        </div>
    );
}

Table.propTypes = {
    columns: PropTypes.array.isRequired,
    data: PropTypes.array.isRequired
}

export default Table; 