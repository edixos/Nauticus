import { createColumnHelper } from '@tanstack/react-table';
import Table from '../../components/Table';
import spacesData from '../../space.json';

const columnHelper = createColumnHelper();

const columns = [
    columnHelper.accessor(row => row.metadata.name, {
        header: 'Name',
        cell: info => info.getValue(),
    }),
    columnHelper.accessor(row => row.spec.resourceQuota.hard['limits.cpu'], {
        header: 'CPU Limits',
        cell: info => info.getValue(),
    }),
    columnHelper.accessor(row => row.spec.resourceQuota.hard['requests.cpu'], {
        header: 'CPU Request',
        cell: info => info.getValue(),
    }),
    columnHelper.accessor(row => row.spec.resourceQuota.hard['limits.memory'], {
        header: 'Memory Limits',
        cell: info => info.getValue(),
    }),
    columnHelper.accessor(row => row.spec.resourceQuota.hard['requests.memory'], {
        header: 'Memory Requests',
        cell: info => info.getValue(),
    }),
    // ... other columns as needed
];


function SpacesPage() {

    return (
        <Table columns={columns} data={spacesData} />
    );

}

export default SpacesPage;