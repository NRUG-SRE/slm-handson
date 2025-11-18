// Playwright MCP動作確認スクリプト
const { chromium } = require('playwright');

(async () => {
  console.log('🚀 Playwright MCP動作テスト開始');

  const browser = await chromium.launch({
    headless: true
  });

  const context = await browser.newContext({
    viewport: { width: 1280, height: 720 }
  });

  const page = await context.newPage();

  console.log('📍 localhost:3000にアクセス中...');
  await page.goto('http://localhost:3000');

  console.log('📸 スクリーンショット撮影中...');
  await page.screenshot({
    path: '.mcp/output/localhost-3000-screenshot.png',
    fullPage: true
  });

  console.log('📄 ページタイトル取得中...');
  const title = await page.title();
  console.log(`   タイトル: ${title}`);

  console.log('🔍 ページ内容確認中...');
  const heading = await page.locator('h1').first().textContent();
  console.log(`   メインヘッダー: ${heading}`);

  const description = await page.locator('p.text-lg').first().textContent();
  console.log(`   説明文: ${description}`);

  // New Relic RUMスクリプトの存在確認
  console.log('🔬 New Relic RUM確認中...');
  const nrScript = await page.evaluate(() => {
    return typeof window.NREUM !== 'undefined';
  });
  console.log(`   New Relic RUM読み込み: ${nrScript ? '✅ 成功' : '❌ 失敗'}`);

  await browser.close();

  console.log('✅ Playwright MCP動作テスト完了');
  console.log(`📸 スクリーンショット保存先: .mcp/output/localhost-3000-screenshot.png`);
})();
