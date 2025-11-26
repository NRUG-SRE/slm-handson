import './style.css'
import Prism from 'prismjs'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-sql'
import 'prismjs/themes/prism.css'

// フルスクリーン状態をグローバルに管理
window.isFullscreenMode = false

// Prism.jsの自動ハイライトを初期化
document.addEventListener('DOMContentLoaded', () => {
  Prism.highlightAll()
  initSlideNavigation()
  initFullscreenMode()
})

// スライドナビゲーション機能の初期化
function initSlideNavigation() {
  const slides = document.querySelectorAll('.slide')
  const prevBtn = document.getElementById('prevSlideBtn')
  const nextBtn = document.getElementById('nextSlideBtn')
  const counter = document.getElementById('slideCounter')

  if (!slides.length || !prevBtn || !nextBtn || !counter) {
    return // スライド要素が存在しない場合は何もしない
  }

  let currentSlide = 0
  const totalSlides = slides.length

  function showSlide(index) {
    // インデックスの範囲チェック
    if (index < 0 || index >= totalSlides) {
      return
    }

    // 全スライドを非表示にし、指定されたスライドのみ表示
    slides.forEach((slide, i) => {
      if (i === index) {
        slide.classList.remove('hidden')
      } else {
        slide.classList.add('hidden')
      }
    })

    // カウンター更新
    counter.textContent = `${index + 1} / ${totalSlides}`

    // ボタンの有効/無効状態を更新
    prevBtn.disabled = index === 0
    nextBtn.disabled = index === totalSlides - 1

    // 現在のスライドインデックスを更新
    currentSlide = index
  }

  // 前へボタンのクリックイベント
  prevBtn.addEventListener('click', () => {
    if (currentSlide > 0) {
      showSlide(currentSlide - 1)
    }
  })

  // 次へボタンのクリックイベント
  nextBtn.addEventListener('click', () => {
    if (currentSlide < totalSlides - 1) {
      showSlide(currentSlide + 1)
    }
  })

  // キーボードナビゲーション（矢印キー）
  document.addEventListener('keydown', (e) => {
    // フルスクリーンモード中はスキップ（フルスクリーン側で処理）
    if (window.isFullscreenMode) return

    // スライドセクションが表示されているかチェック
    const lectureSection = document.getElementById('lecture')
    if (!lectureSection) return

    // 座学セクションが画面内にあるかチェック（簡易的な判定）
    const rect = lectureSection.getBoundingClientRect()
    const isVisible = rect.top < window.innerHeight && rect.bottom > 0

    if (!isVisible) return

    if (e.key === 'ArrowLeft' && currentSlide > 0) {
      e.preventDefault()
      showSlide(currentSlide - 1)
    } else if (e.key === 'ArrowRight' && currentSlide < totalSlides - 1) {
      e.preventDefault()
      showSlide(currentSlide + 1)
    }
  })

  // 初期表示
  showSlide(0)

  // 現在のスライドインデックスを外部から取得できるようにする
  window.getCurrentSlide = () => currentSlide
  window.setCurrentSlide = (index) => showSlide(index)
  window.getTotalSlides = () => totalSlides
}

// フルスクリーンモード機能の初期化
function initFullscreenMode() {
  const fullscreenBtn = document.getElementById('fullscreenBtn')
  const exitFullscreenBtn = document.getElementById('exitFullscreenBtn')
  const fullscreenModal = document.getElementById('fullscreenModal')
  const fullscreenContent = document.getElementById('fullscreenContent')
  const fullscreenPrevBtn = document.getElementById('fullscreenPrevBtn')
  const fullscreenNextBtn = document.getElementById('fullscreenNextBtn')
  const fullscreenCounter = document.getElementById('fullscreenCounter')
  const slides = document.querySelectorAll('.slide')

  if (!fullscreenBtn || !fullscreenModal || !slides.length) {
    return
  }

  // フルスクリーンでスライドを表示
  function showFullscreenSlide(index) {
    const totalSlides = window.getTotalSlides()
    if (index < 0 || index >= totalSlides) return

    // 現在のスライドをフルスクリーンコンテンツにコピー
    const currentSlideElement = slides[index]
    fullscreenContent.innerHTML = ''

    // スライドのクローンを作成してスタイルを調整
    const clone = currentSlideElement.cloneNode(true)
    clone.classList.remove('hidden')
    clone.classList.add('w-full', 'max-w-6xl', 'text-xl')

    // フォントサイズを拡大
    clone.querySelectorAll('h3').forEach(el => {
      el.classList.remove('text-2xl', 'text-3xl')
      el.classList.add('text-4xl', 'mb-8')
    })
    clone.querySelectorAll('p').forEach(el => {
      el.classList.add('text-xl')
    })
    clone.querySelectorAll('li').forEach(el => {
      el.classList.add('text-lg')
    })

    fullscreenContent.appendChild(clone)

    // カウンター更新
    fullscreenCounter.textContent = `${index + 1} / ${totalSlides}`

    // ボタンの有効/無効状態を更新
    fullscreenPrevBtn.disabled = index === 0
    fullscreenNextBtn.disabled = index === totalSlides - 1

    // メインのスライドも同期
    window.setCurrentSlide(index)
  }

  // フルスクリーンモードを開く
  function openFullscreen() {
    window.isFullscreenMode = true
    fullscreenModal.classList.remove('hidden')
    fullscreenModal.classList.add('flex')
    document.body.style.overflow = 'hidden'
    showFullscreenSlide(window.getCurrentSlide())
  }

  // フルスクリーンモードを閉じる
  function closeFullscreen() {
    window.isFullscreenMode = false
    fullscreenModal.classList.add('hidden')
    fullscreenModal.classList.remove('flex')
    document.body.style.overflow = ''
  }

  // イベントリスナー
  fullscreenBtn.addEventListener('click', openFullscreen)
  exitFullscreenBtn.addEventListener('click', closeFullscreen)

  // フルスクリーン内のナビゲーション
  fullscreenPrevBtn.addEventListener('click', () => {
    const current = window.getCurrentSlide()
    if (current > 0) {
      showFullscreenSlide(current - 1)
    }
  })

  fullscreenNextBtn.addEventListener('click', () => {
    const current = window.getCurrentSlide()
    const total = window.getTotalSlides()
    if (current < total - 1) {
      showFullscreenSlide(current + 1)
    }
  })

  // ESCキーで閉じる & 矢印キーでナビゲーション
  document.addEventListener('keydown', (e) => {
    if (!window.isFullscreenMode) return

    if (e.key === 'Escape') {
      closeFullscreen()
    } else if (e.key === 'ArrowLeft') {
      e.preventDefault()
      const current = window.getCurrentSlide()
      if (current > 0) {
        showFullscreenSlide(current - 1)
      }
    } else if (e.key === 'ArrowRight') {
      e.preventDefault()
      const current = window.getCurrentSlide()
      const total = window.getTotalSlides()
      if (current < total - 1) {
        showFullscreenSlide(current + 1)
      }
    }
  })
}
