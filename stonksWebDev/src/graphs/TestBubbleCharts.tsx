import {
  Chart as ChartJS,
  LinearScale,
  PointElement,
  Tooltip,
  Legend,
} from "chart.js";
import { Bubble } from "react-chartjs-2";

ChartJS.register(LinearScale, PointElement, Tooltip, Legend);

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
      label: "Red dataset",
      data: [
        [2, 2, 1],
        [4, 3, 10],
        [5, 40, 3],
        [6, 5, 5],
        [10, 4, 20],
      ],
      backgroundColor: "rgba(255, 99, 132, 0.5)",
    },
    {
      label: "Blue dataset",
      data: [
        [0, 2, 3],
        [1, 4, 2],
        [2, 10, 5],
        [4, 2, 30],
        [20, 3, 12],
      ],
      backgroundColor: "rgba(53, 162, 235, 0.5)",
    },
  ],
};

export function BubbleChartsTest() {
  return <Bubble options={options} data={data} />;
}
