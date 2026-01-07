import { createFileRoute, redirect } from "@tanstack/react-router";
import { isAuthenticated } from "../../lib/auth";
import CustomPage from "../../pages/CustomPage";

export const Route = createFileRoute("/pages/$page")({
  beforeLoad: ({ location }: any) => {
    if (!isAuthenticated()) {
      throw redirect({
        to: "/login",
        search: { redirect: location.href },
      });
    }
  },
  component: CustomPage,
});
