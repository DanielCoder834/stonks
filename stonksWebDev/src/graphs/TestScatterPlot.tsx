import {
  Chart as ChartJS,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
} from "chart.js";
import { Scatter } from "react-chartjs-2";

ChartJS.register(LinearScale, PointElement, LineElement, Tooltip, Legend);

export const options = {
  scales: {
    y: {
      beginAtZero: true,
    },
  },
};

export const data = {
  datasets: [
    {
      label: "A dataset",
      data: [
        [0.2, 0.1],
        [0.3, 0.1],
        [0.3, 0.2],
        [0.2, 1],
        [0.4, 2],
        [10, 2],
        [5, 1],
        [3, 4],
        [1, 3],
      ],
      backgroundColor: "rgba(255, 99, 132, 1)",
    },
  ],
};

export function TestScatter() {
  return <Scatter options={options} data={data} />;
}
