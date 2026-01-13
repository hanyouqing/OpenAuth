import { message } from 'antd'
import { AxiosError } from 'axios'

export function useErrorHandler() {
  const handleError = (error: unknown) => {
    if (error instanceof AxiosError) {
      const errorMessage =
        error.response?.data?.error ||
        error.message ||
        'An error occurred'
      message.error(errorMessage)
    } else if (error instanceof Error) {
      message.error(error.message)
    } else {
      message.error('An unexpected error occurred')
    }
  }

  return { handleError }
}
