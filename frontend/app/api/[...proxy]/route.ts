import { NextRequest, NextResponse } from 'next/server'

// Docker内部通信用のAPI URL
const INTERNAL_API_URL = process.env.INTERNAL_API_URL || 'http://api-server:8080/api'

export async function GET(
  request: NextRequest,
  { params }: { params: { proxy: string[] } }
) {
  const apiPath = params.proxy.join('/')
  const searchParams = request.nextUrl.searchParams.toString()
  const queryString = searchParams ? `?${searchParams}` : ''
  
  try {
    // フロントエンドから送られたヘッダーをバックエンドに転送
    const forwardHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    
    // Distributed Tracing と Session ID のヘッダーを転送
    const newrelicHeader = request.headers.get('newrelic')
    const sessionIdHeader = request.headers.get('X-Session-ID')
    const newrelicTraceHeader = request.headers.get('X-NewRelic-Trace')
    const newrelicBrowserHeader = request.headers.get('X-NewRelic-Browser')
    
    if (newrelicHeader) {
      forwardHeaders['newrelic'] = newrelicHeader
    }
    if (sessionIdHeader) {
      forwardHeaders['X-Session-ID'] = sessionIdHeader
    }
    if (newrelicTraceHeader) {
      forwardHeaders['X-NewRelic-Trace'] = newrelicTraceHeader
    }
    if (newrelicBrowserHeader) {
      forwardHeaders['X-NewRelic-Browser'] = newrelicBrowserHeader
    }
    
    const response = await fetch(`${INTERNAL_API_URL}/${apiPath}${queryString}`, {
      method: 'GET',
      headers: forwardHeaders,
    })

    if (!response.ok) {
      return NextResponse.json(
        { error: 'API request failed' },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('API proxy error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function POST(
  request: NextRequest,
  { params }: { params: { proxy: string[] } }
) {
  const apiPath = params.proxy.join('/')
  
  try {
    const body = await request.text()
    
    // フロントエンドから送られたヘッダーをバックエンドに転送
    const forwardHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    
    // Distributed Tracing と Session ID のヘッダーを転送
    const newrelicHeader = request.headers.get('newrelic')
    const sessionIdHeader = request.headers.get('X-Session-ID')
    const newrelicTraceHeader = request.headers.get('X-NewRelic-Trace')
    const newrelicBrowserHeader = request.headers.get('X-NewRelic-Browser')
    
    if (newrelicHeader) {
      forwardHeaders['newrelic'] = newrelicHeader
    }
    if (sessionIdHeader) {
      forwardHeaders['X-Session-ID'] = sessionIdHeader
    }
    if (newrelicTraceHeader) {
      forwardHeaders['X-NewRelic-Trace'] = newrelicTraceHeader
    }
    if (newrelicBrowserHeader) {
      forwardHeaders['X-NewRelic-Browser'] = newrelicBrowserHeader
    }
    
    const response = await fetch(`${INTERNAL_API_URL}/${apiPath}`, {
      method: 'POST',
      headers: forwardHeaders,
      body,
    })

    if (!response.ok) {
      return NextResponse.json(
        { error: 'API request failed' },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('API proxy error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function PUT(
  request: NextRequest,
  { params }: { params: { proxy: string[] } }
) {
  const apiPath = params.proxy.join('/')
  
  try {
    const body = await request.text()
    
    // フロントエンドから送られたヘッダーをバックエンドに転送
    const forwardHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    
    // Distributed Tracing と Session ID のヘッダーを転送
    const newrelicHeader = request.headers.get('newrelic')
    const sessionIdHeader = request.headers.get('X-Session-ID')
    const newrelicTraceHeader = request.headers.get('X-NewRelic-Trace')
    const newrelicBrowserHeader = request.headers.get('X-NewRelic-Browser')
    
    if (newrelicHeader) {
      forwardHeaders['newrelic'] = newrelicHeader
    }
    if (sessionIdHeader) {
      forwardHeaders['X-Session-ID'] = sessionIdHeader
    }
    if (newrelicTraceHeader) {
      forwardHeaders['X-NewRelic-Trace'] = newrelicTraceHeader
    }
    if (newrelicBrowserHeader) {
      forwardHeaders['X-NewRelic-Browser'] = newrelicBrowserHeader
    }
    
    const response = await fetch(`${INTERNAL_API_URL}/${apiPath}`, {
      method: 'PUT',
      headers: forwardHeaders,
      body,
    })

    if (!response.ok) {
      return NextResponse.json(
        { error: 'API request failed' },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('API proxy error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}