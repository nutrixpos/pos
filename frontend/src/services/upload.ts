import axios, { AxiosProgressEvent } from 'axios'


const upload_svc:any = {}

type ProgressEventHandler = (progressEvent: AxiosProgressEvent) => void;


upload_svc.uploadProductImage = (product_id: String,file: File,backendURL: string,progressLoadingFunc: ProgressEventHandler) =>  {
    

        let formData = new FormData();
        formData.append('file', file);
        formData.append('type', file.type);

        return axios.patch(`${backendURL}/api/products/${product_id}/image`,
            formData,
            {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
                onUploadProgress: progressLoadingFunc
                // onUploadProgress: progressEvent => {
                //     // uploadingProgressText = `Uploading media ${i+1}/${media.value.length} (${Math.floor(progressEvent.loaded * 100 / size)}%)`
                //     if (mediaLoading.value.length < i+1)
                //         mediaLoading.value.push({
                //             isLoading : true
                //         } as MediaLoading)

                //     mediaLoading.value[i].isLoading = true
                //     mediaLoading.value[i].percent = Math.floor(progressEvent.loaded * 100 / size)
                // }
            }
        )
}



export default upload_svc