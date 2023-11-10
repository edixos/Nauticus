import spacesData from '../../space.json';
import TableBody from './components/TableBody';
import TableHeader from './components/TableHeader';

function Spaces() {
    console.log(spacesData);
    const tableHeader = [
        { name: 'Name' },
        { name: 'Owners' },
        { name: 'CPU Limit' },
        { name: 'Memory Limit' },
    ]

    return (
        <div className="overflow-x-auto p-1">
            <table className="min-w-full table-auto bg-slate-800 shadow-md rounded-md overflow-hidden">
                <TableHeader headers={tableHeader} />

                <TableBody data={spacesData} />
            </table>
        </div>
    );
}

export default Spaces;