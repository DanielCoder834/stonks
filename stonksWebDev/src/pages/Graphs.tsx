import { TestGraph } from "../graphs/TestLineGraph";
import { TestPieChart } from "../graphs/TestPieCharts";
import { VertBarGraph } from "../graphs/TestVeritcalBarGraph";
import { DoughnutTest } from "../graphs/TestDoughnutCharts";
import { TestScatter } from "../graphs/TestScatterPlot";
import { BubbleChartsTest } from "../graphs/TestBubbleCharts";

export default function Graphs() {
  return (
    <>
      <BubbleChartsTest />
      <TestScatter />
      <DoughnutTest />
      <TestPieChart />
      <TestGraph />
      <VertBarGraph />
    </>
  );
}
