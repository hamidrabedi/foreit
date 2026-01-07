import { useEffect, useRef } from 'react';
import { useToast } from './use-toast';
import { getAdminConfig } from '../lib/config';

export function useNotifications() {
  const { toast } = useToast();
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    // Determine SSE URL (relative to API base)
    const { apiBase } = getAdminConfig();
    const protocol = window.location.protocol;
    const host = window.location.host;
    const sseUrl = `${protocol}//${host}${apiBase}/events`;

    console.log('[Notifications] Connecting to SSE:', sseUrl);
    
    const es = new EventSource(sseUrl);
    eventSourceRef.current = es;

    es.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('[Notifications] New event:', data);
        
        toast({
          title: data.title || 'Notification',
          description: data.message,
          variant: data.type === 'error' ? 'destructive' : 'default',
        });
      } catch (err) {
        console.error('[Notifications] Failed to parse event data:', err);
      }
    };

    es.onerror = (err) => {
      console.error('[Notifications] SSE Error:', err);
      // EventSource automatically retries, but we might want to log it
    };

    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    };
  }, [toast]);
}
