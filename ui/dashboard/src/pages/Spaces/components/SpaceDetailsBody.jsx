import Label from '../../../components/Label';
import { HiMiniCpuChip } from "react-icons/hi2";
import { FaMemory } from "react-icons/fa6";
import { PropTypes } from 'prop-types';
import SpaceDetailsBodySection from './SpaceDetailsBodySection';


function SpaceDetailsBody({ spaceDetails }) {

    // First, safely access nested properties
    const cpuLimits = spaceDetails?.spec?.resourceQuota?.hard?.["limits.cpu"];
    const memoryLimits = spaceDetails?.spec?.resourceQuota?.hard?.["limits.memory"];
    // Assuming you also have usage data
    const cpuUsage = spaceDetails?.status?.resourceQuota?.used?.cpu;
    const memoryUsage = spaceDetails?.status?.resourceQuota?.used?.memory;

    // Calculate width percentages for progress bars
    const cpuWidthPercent = `${(parseFloat(cpuUsage) / parseFloat(cpuLimits)) * 100}%`;
    const memoryWidthPercent = `${(parseFloat(memoryUsage) / parseFloat(memoryLimits)) * 100}%`;

    console.log(memoryWidthPercent);
    return (
        <div className="border-t border-gray-200">
            <dl>

                {/* Labels */}
                <SpaceDetailsBodySection title="Labels">
                    {Object.entries(spaceDetails.metadata.labels).map(([key, value]) => (
                        <Label key={key} className="bg-green-400 text-white" >
                            {key}: {value}
                        </Label>
                    ))}
                </SpaceDetailsBodySection>

                {/* Owner */}
                <SpaceDetailsBodySection title="Owners">
                    {spaceDetails.spec.owners.map(owner => (
                        <div key={owner.name} className="flex items-center">
                            <div className="text-sm font-medium text-gray-900">
                                {owner.name}
                            </div>
                            <div className="ml-4 text-sm text-gray-500">
                                ({owner.kind})
                            </div>
                        </div>
                    ))}
                </SpaceDetailsBodySection>

                <SpaceDetailsBodySection title="Capacity">
                    {/* CPU */}
                    <div className="flex items-center space-x-3">
                        <HiMiniCpuChip className="text-2xl" /> {/* CPU Icon */}
                        <div className="flex-1 bg-gray-200 rounded-full h-6 overflow-hidden">
                            <div className="bg-blue-600 h-6 rounded-full" style={{ width: cpuWidthPercent }}>
                                <span className="block text-white text-sm text-center">{cpuUsage} / {cpuLimits}</span>
                            </div>
                        </div>
                    </div>
                    {/* Memory */}
                    <div className="flex items-center space-x-3 mt-2">
                        <FaMemory className="text-2xl" /> {/* Memory Icon */}
                        <div className="flex-1 bg-gray-200 rounded-full h-6 overflow-hidden">
                            <div className="bg-blue-600 h-6 rounded-full" style={{ width: memoryWidthPercent }}>
                                <span className="block text-white text-sm text-center">{memoryUsage} / {memoryLimits}</span>
                            </div>
                        </div>
                    </div>
                </SpaceDetailsBodySection>
            </dl>
        </div>
    );
}

SpaceDetailsBody.propTypes = {
    spaceDetails: PropTypes.object,
}


export default SpaceDetailsBody;