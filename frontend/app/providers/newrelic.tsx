'use client';

import { useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { monitoring } from '@/lib/monitoring';

export default function NewRelicRouteTracker() {
  const pathname = usePathname();

  useEffect(() => {
    // New Relic監視の初期化
    monitoring.init();

    if (typeof window !== 'undefined' && window.newrelic) {
      // 仮想ページ遷移時にNew Relicへ通知
      if (window.newrelic.setPageViewName) {
        window.newrelic.setPageViewName(pathname);
      }
      // より詳細なページアクションも記録
      window.newrelic.addPageAction('routeChange', {
        pathname: pathname,
        timestamp: Date.now(),
        userAgent: navigator.userAgent
      });
    }
  }, [pathname]);

  return null;
}