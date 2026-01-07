import { createFileRoute, redirect } from '@tanstack/react-router';
import DashboardPage from '../pages/DashboardPage';
import { isAuthenticated } from '../lib/auth';

export const Route = createFileRoute('/')({
  beforeLoad: ({ location }: any) => {
    if (!isAuthenticated()) {
      throw redirect({
        to: '/login',
        search: { redirect: location.href },
      });
    }
  },
  component: DashboardPage,
});
