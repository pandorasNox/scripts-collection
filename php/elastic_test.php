<?php

echo PHP_EOL;
echo "start script";
echo PHP_EOL;


# composer vendor
require_once('./../vendor/autoload.php');

use Elasticsearch\ClientBuilder;

$hosts = [
    [
        'host' => 'elasticsearch',
        'port' => '9200',
        // 'scheme' => 'https',
        'user' => 'elastic',
        'pass' => 'changeme'
    ],
];
$client = ClientBuilder::create()           // Instantiate a new ClientBuilder
            ->setHosts($hosts)      // Set the hosts
            ->build();              // Build the client object

try {
    #$status = $client->getStatus();
    $response = $client->indices()->stats();
    var_dump($response);
} catch(Exception $e) {
    var_dump($e);
}



echo PHP_EOL;
echo "script finished";
echo PHP_EOL;
