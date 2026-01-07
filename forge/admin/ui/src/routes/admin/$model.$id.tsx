import { createFileRoute, redirect } from '@tanstack/react-router';
import DynamicModelPage from '../../components/DynamicModelPage';
import { isAuthenticated } from '../../lib/auth';

export const Route = createFileRoute('/$model/$id')({
  beforeLoad: ({ location }: any) => {
    if (!isAuthenticated()) {
      throw redirect({
        to: '/login',
        search: { redirect: location.href },
      });
    }
  },
  component: () => <DynamicModelPage mode="edit" />,
});
