import { createFileRoute } from "@tanstack/react-router";
import DynamicModelPage from "../../components/DynamicModelPage";

export const Route = createFileRoute("/admin/$model/new")({
  component: () => <DynamicModelPage mode="create" />,
});
