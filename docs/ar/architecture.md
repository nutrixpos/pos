---
outline: deep
---

# الهندسة المعمارية

Nutrix هو تطبيق ويب يتكون من الأنظمة الفرعية الرئيسية الثلاثة التالية:

#### pos (الواجهة الخلفية)
النظام الأساسي والخادم الخلفي لـ Nutrix ، وهي واجهة برمجة تطبيقات REST مبنية في [golang](go.dev) ، لمزيد من المعلومات تحقق من مستودع git المخصص [nutrixpos/pos](https://github.com/nutrixpos/pos) ، يستخدم pos قاعدة بيانات قائمة على المستندات المدمجة بـ [mongo](https://www.mongodb.com/)

#### posui (الواجهة الأمامية)
واجهة المستخدم لـ Nutrix ، تطبيق ويب مبني في [vue](https://vuejs.org/) ، يتم تضمين الكود في دليل فرعي [pos/frontend](https://github.com/nutrixpos/pos/tree/main/frontend).