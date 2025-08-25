const { chromium } = require('playwright-chromium');

class UltraLoadTester {
  constructor() {
    this.config = {
      targetUrl: process.env.TARGET_URL || 'http://frontend:3000',
      accessInterval: parseInt(process.env.ACCESS_INTERVAL || '10') * 1000,
      duration: parseInt(process.env.DURATION || '300') * 1000,
      headless: process.env.HEADLESS !== 'false',
      concurrentUsers: parseInt(process.env.CONCURRENT_USERS || '1'),
    };

    this.stats = {
      total: 0,
      success: 0,
      failed: 0,
      times: [],
      startTime: Date.now()
    };
  }

  log(msg) {
    const time = new Date().toISOString().split('T')[1].split('.')[0];
    console.log(`${time} ${msg}`);
  }

  sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async userJourney(userId) {
    const browser = await chromium.launch({
      headless: this.config.headless,
      executablePath: '/usr/bin/chromium',
      args: [
        '--no-sandbox',
        '--disable-setuid-sandbox',
        '--disable-dev-shm-usage',
        '--disable-gpu',
        '--single-process',
        '--disable-web-security',
        '--disable-images',
        '--memory-pressure-off',
        '--disable-extensions',
        '--disable-plugins'
      ]
    });

    const context = await browser.newContext({
      userAgent: `SLM-LoadTester-${userId}`,
      viewport: { width: 1280, height: 720 }
    });

    const endTime = Date.now() + this.config.duration;

    try {
      while (Date.now() < endTime) {
        const page = await context.newPage();
        const journeyStart = Date.now();

        try {
          // 1. TOPページアクセス
          const topPageStart = Date.now();
          this.log(`顧客${userId}: 🏠 ECサイトにアクセス`);
          const response = await page.goto(this.config.targetUrl, {
            waitUntil: 'domcontentloaded',
            timeout: 20000
          });
          const topPageTime = Date.now() - topPageStart;

          if (response && response.ok()) {
            this.log(`顧客${userId}: 📱 商品一覧を閲覧中... (読み込み: ${topPageTime}ms)`);
            
            // 2. 商品詳細ページへ遷移
            await this.sleep(2000 + Math.random() * 3000); // 2-5秒の閲覧時間
            
            const links = await page.$$('a[href*="/products/"]').catch(() => []);
            if (links.length > 0) {
              const randomIndex = Math.floor(Math.random() * links.length);
              const randomLink = links[randomIndex];
              
              // 商品名を取得（実際のページから）
              await page.waitForSelector('h1', { timeout: 10000 }).catch(() => {});
              const productPageStart = Date.now();
              await randomLink.click().catch(() => {});
              await page.waitForLoadState('domcontentloaded').catch(() => {});
              
              // React/Next.jsの状態が安定するまで待つ
              await this.sleep(1500);
              await page.waitForSelector('h1', { timeout: 10000 }).catch(() => {});
              
              // 商品名を取得
              const productName = await page.$eval('h1', el => el.textContent).catch(() => '商品');
              const productPageTime = Date.now() - productPageStart;
              
              this.log(`顧客${userId}: 🛍️ 「${productName}」の詳細を確認中... (読み込み: ${productPageTime}ms)`);
              
              // 3. 商品詳細の確認とカート操作
              await this.sleep(3000 + Math.random() * 4000); // 3-7秒の商品検討時間
              
              // カート追加ボタンを探す
              const cartActionStart = Date.now();
              const addButton = await page.$('button:has-text("カートに追加")').catch(() => null);
              
              if (addButton) {
                // ボタンが有効（disabled状態でない）かチェック
                const isDisabled = await addButton.getAttribute('disabled');
                if (!isDisabled) {
                  this.log(`顧客${userId}: 🛒 「${productName}」をカートに追加...`);
                  await addButton.click().catch(() => {});
                  await this.sleep(800); // カート追加後の待機
                  
                  // カート追加完了メッセージが表示されるまで待機
                  await page.waitForSelector('.bg-green-100, .bg-red-100', { timeout: 3000 }).catch(() => {});
                  const cartActionTime = Date.now() - cartActionStart;
                  this.log(`顧客${userId}: ✅ カートに追加完了！ (処理時間: ${cartActionTime}ms)`);
                } else {
                  this.log(`顧客${userId}: 😞 「${productName}」は在庫切れでした`);
                }
              } else {
                this.log(`顧客${userId}: ⚠️ カート機能に問題が発生（SLO違反の可能性）`);
              }

              // 4. カートページ表示
              await this.sleep(1000 + Math.random() * 2000); // カート確認前の思考時間
              const cartPageStart = Date.now();
              this.log(`顧客${userId}: 💳 カートの内容を確認中...`);
              const cartResponse = await page.goto(`${this.config.targetUrl}/cart`, {
                waitUntil: 'domcontentloaded',
                timeout: 15000
              }).catch(() => null);
              const cartPageTime = Date.now() - cartPageStart;
              
              if (cartResponse && cartResponse.ok()) {
                this.log(`顧客${userId}: ✅ カート内容を確認完了 (表示時間: ${cartPageTime}ms)`);
                
                // 5. カートページでの確認時間
                await this.sleep(2000 + Math.random() * 3000); // 2-5秒のカート確認時間
                
                // 6. 決済ページへ遷移
                const checkoutPageStart = Date.now();
                this.log(`顧客${userId}: 💳 決済手続きを開始...`);
                const checkoutResponse = await page.goto(`${this.config.targetUrl}/checkout`, {
                  waitUntil: 'domcontentloaded',
                  timeout: 15000
                }).catch(() => null);
                const checkoutPageTime = Date.now() - checkoutPageStart;
                
                if (checkoutResponse && checkoutResponse.ok()) {
                  this.log(`顧客${userId}: 📝 決済ページ表示完了 (読み込み: ${checkoutPageTime}ms)`);
                  
                  // 7. 決済情報入力・確認時間
                  await this.sleep(3000 + Math.random() * 5000); // 3-8秒の決済検討時間
                  
                  // 8. 注文確定処理
                  const orderStart = Date.now();
                  this.log(`顧客${userId}: 🎯 注文を確定中...`);
                  
                  // 注文確定ボタンを探してクリック
                  const orderButton = await page.$('button:has-text("注文を確定"), button:has-text("決済"), button[type="submit"]').catch(() => null);
                  
                  if (orderButton) {
                    const isDisabled = await orderButton.getAttribute('disabled');
                    if (!isDisabled) {
                      await orderButton.click().catch(() => {});
                      await this.sleep(1000 + Math.random() * 2000); // 決済処理待機
                      
                      // 成功メッセージまたはページ変更を待つ
                      await page.waitForSelector('.bg-green-100, .bg-blue-100, h1', { timeout: 5000 }).catch(() => {});
                      
                      const orderTime = Date.now() - orderStart;
                      this.log(`顧客${userId}: 🎉 購入完了！ (決済処理: ${orderTime}ms)`);
                    } else {
                      this.log(`顧客${userId}: ⚠️ 決済ボタンが無効状態`);
                    }
                  } else {
                    this.log(`顧客${userId}: ❌ 決済ボタンが見つかりません (UX問題)`);
                  }
                } else {
                  this.log(`顧客${userId}: ❌ 決済ページでエラー発生 (SLO違反: ${checkoutPageTime}ms)`);
                }
              } else {
                this.log(`顧客${userId}: ❌ カートページでエラー発生 (SLO違反: ${cartPageTime}ms)`);
              }
            } else {
              this.log(`顧客${userId}: 😵 商品が表示されません（重大なエラー）`);
            }

            const totalTime = Date.now() - journeyStart;
            this.stats.times.push(totalTime);
            this.stats.success++;
            this.log(`顧客${userId}: ✨ 購入体験完了 (所要時間: ${Math.round(totalTime/1000)}秒) - New Relicでコンバージョン確認可能`);
          } else {
            this.stats.failed++;
            this.log(`顧客${userId}: ❌ サイトにアクセスできません (${response ? response.status() : '接続失敗'}) - SLO違反`);
          }
        } catch (error) {
          this.stats.failed++;
          this.log(`顧客${userId}: ⚠️ システムエラーが発生しました - ${error.message}`);
        } finally {
          this.stats.total++;
          await page.close().catch(() => {});
        }

        // 次の顧客のアクセスまで待機
        if (Date.now() < endTime) {
          const baseInterval = this.config.accessInterval;
          const randomVariation = Math.random() * baseInterval * 0.5;
          const actualInterval = baseInterval + (randomVariation - baseInterval * 0.25);
          this.log(`💤 次の顧客アクセスまで ${Math.round(actualInterval/1000)}秒後...`);
          await this.sleep(actualInterval);
        }
      }
    } finally {
      await context.close().catch(() => {});
      await browser.close().catch(() => {});
    }
  }

  printStats() {
    const avg = this.stats.times.length > 0 
      ? Math.round(this.stats.times.reduce((a, b) => a + b, 0) / this.stats.times.length)
      : 0;
    
    const sorted = [...this.stats.times].sort((a, b) => a - b);
    const p95 = sorted[Math.floor(sorted.length * 0.95)] || 0;
    const p50 = sorted[Math.floor(sorted.length * 0.5)] || 0;
    const min = sorted[0] || 0;
    const max = sorted[sorted.length - 1] || 0;
    
    const elapsed = Math.round((Date.now() - this.stats.startTime) / 1000);
    const successRate = this.stats.total > 0 ? ((this.stats.success / this.stats.total) * 100).toFixed(1) : '0.0';
    const throughput = this.stats.total > 0 ? (this.stats.total / elapsed * 60).toFixed(1) : '0.0';

    console.log('\n' + '='.repeat(50));
    console.log('📊 SLO/SLI監視データ生成状況');
    console.log('='.repeat(50));
    console.log(`⏱️  稼働時間: ${elapsed}秒 / 予定: ${Math.round(this.config.duration/1000)}秒`);
    console.log(`👥 アクセス実行: ${this.stats.total}回 (成功: ${this.stats.success}回, 失敗: ${this.stats.failed}回)`);
    console.log('');
    console.log('📈 アプリケーションパフォーマンス:');
    console.log(`   レスポンス時間 - 最小: ${min}ms | 平均: ${avg}ms | P95: ${p95}ms`);
    if (this.stats.failed > 0) {
      console.log(`   ⚠️  エラー発生: ${this.stats.failed}件 - SLO違反の可能性`);
    }
    console.log('='.repeat(50));
    console.log('💡 New Relic UIでリアルタイムSLI/SLO監視データを確認');
    console.log('='.repeat(50));
  }

  async run() {
    this.log('🚀 SLMハンズオン - バーチャル顧客体験シミュレーション開始');
    this.log(`📊 設定: ${this.config.concurrentUsers}名の顧客が ${this.config.duration / 1000}秒間利用`);
    this.log(`🌐 対象ECサイト: ${this.config.targetUrl}`);
    this.log('📈 New Relic RUM & APMでリアルタイム監視中...');

    const customers = [];
    for (let i = 1; i <= this.config.concurrentUsers; i++) {
      customers.push(this.userJourney(i));
    }

    const statsInterval = setInterval(() => this.printStats(), 30000);
    
    try {
      await Promise.all(customers);
    } finally {
      clearInterval(statsInterval);
      this.printStats();
      this.log('🎯 顧客体験シミュレーション完了！');
      this.log('📊 New Relic UIでSLO/SLI違反やパフォーマンス劣化を確認してください');
      this.log('💡 エラーバジェット消費状況をダッシュボードで監視可能です');
    }
  }
}

// グレースフルシャットダウン
process.on('SIGINT', () => {
  console.log('\n⏹️  停止中...');
  process.exit(0);
});

process.on('SIGTERM', () => {
  console.log('\n⏹️  停止中...');
  process.exit(0);
});

// 実行
(async () => {
  try {
    const tester = new UltraLoadTester();
    await tester.run();
    process.exit(0);
  } catch (error) {
    console.error('🔥 エラー:', error);
    process.exit(1);
  }
})();