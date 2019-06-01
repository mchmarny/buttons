
## Data

BigQuery has now PubSub data source connector (Cloud Dataflow behind the scene) so you can use SQL to query PubSub topic payloads while it's streaming. The full instructions BigQuery are [here](https://cloud.google.com/dataflow/docs/guides/sql/dataflow-sql-ui-walkthrough). Follow the instruction there to enable the API and create the necessary service account.

### Register Schema

```shell
gcloud beta data-catalog entries update \
    --lookup-entry='pubsub.topic.`project_ID`.clicks' \
    --schema-from-file=schema.yaml
```

### Create Job

 And using [Tumbling windows function](https://cloud.google.com/dataflow/docs/guides/sql/streaming-pipeline-basics#tumbling-windows) to display the streamed data in non overlapping time interval.

```sql
SELECT
    c.payload.color as button_color,
    TUMBLE_START("INTERVAL 5 SECOND") AS period_start,
    SUM(c.payload.click) as click_sum
FROM pubsub.topic.`project_ID`.clicks as c
GROUP BY
    c.payload.color,
    TUMBLE(c.event_timestamp, "INTERVAL 5 SECOND")
```

### Select Clicks

```sql
SELECT * FROM buttons.clicks_5s order by period_start desc
```

Should return

| Row	| button_color	| period_start	            | click_sum	|
| 1	    | white         | 2019-05-31 23:16:25 UTC   | 14        |
| 2	    | white         | 2019-05-31 23:13:20 UTC   | 1         |
| 3     | green         | 2019-05-31 23:13:20 UTC   | 1         |
| 4	    | black         | 2019-05-31 23:13:15 UTC   | 4         |

