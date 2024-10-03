export class Notification {
    description: string
    type: string
    topic_name: string
    severity: string
    constructor(){
        this.description = ""
        this.type = ""
        this.topic_name = ""
        this.severity = ""
    }
}


export class OrderFinishNotification extends Notification {
    order_id: string
    constructor(){
        super()
        this.order_id = ""
    }
}