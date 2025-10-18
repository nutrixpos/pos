export class Notification {
    id: string
    date: Date
    description: string
    type: string
    topic_name: string
    severity: string
    constructor(){
        this.id = Math.random().toString(16).substring(2, 34)
        this.description = ""
        this.type = ""
        this.topic_name = ""
        this.severity = ""
        this.date = new Date()
    }
}


export class OrderFinishNotification extends Notification {
    order_id: string
    constructor(){
        super()
        this.order_id = ""
    }
}