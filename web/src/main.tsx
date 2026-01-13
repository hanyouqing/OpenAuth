import React from 'react'
import ReactDOM from 'react-dom/client'
import { ConfigProvider } from 'antd'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import ErrorBoundary from './components/ErrorBoundary'
import App from './App'
import './i18n'
import './styles/index.css'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
})

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <ConfigProvider
          theme={{
            token: {
              colorPrimary: '#000000',
              colorSuccess: '#000000',
              colorWarning: '#616161',
              colorError: '#000000',
              colorInfo: '#212121',
              colorBgBase: '#FFFFFF',
              colorText: '#000000',
              colorTextSecondary: '#616161',
              colorBorder: '#E0E0E0',
              colorBorderSecondary: '#F5F5F5',
              borderRadius: 4,
              fontFamily:
                '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
            },
            components: {
              Button: {
                borderRadius: 4,
                controlHeight: 40,
                primaryColor: '#FFFFFF',
                primaryShadow: 'none',
              },
              Input: {
                borderRadius: 4,
                controlHeight: 40,
                activeBorderColor: '#000000',
                hoverBorderColor: '#212121',
              },
              Card: {
                borderRadius: 4,
                headerBg: '#F5F5F5',
              },
              Menu: {
                itemBorderRadius: 4,
                itemSelectedBg: '#F5F5F5',
                itemHoverBg: '#F5F5F5',
              },
              Layout: {
                bodyBg: '#FFFFFF',
                headerBg: '#FFFFFF',
                siderBg: '#F5F5F5',
              },
            },
          }}
        >
          <BrowserRouter
            future={{
              v7_startTransition: true,
              v7_relativeSplatPath: true,
            }}
          >
            <App />
          </BrowserRouter>
        </ConfigProvider>
      </QueryClientProvider>
    </ErrorBoundary>
  </React.StrictMode>,
)
