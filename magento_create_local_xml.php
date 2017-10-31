<?php

$pathArg = array_search("--path", $argv);

if (false === $pathArg) {
    fwrite(STDERR, "An error occurred. No path provided. Please add a (relative) path with '--path ./../path/to/magento/app/etc/'." . PHP_EOL);
    exit(1); // A response code other than 0 is a failure
}

if (false === isset($argv[$pathArg + 1])) {
    fwrite(STDERR, "An error occurred. No path provided after '--path'. Please add a (relative) path with '--path ./../path/to/magento/app/etc/'." . PHP_EOL);
    exit(1); // A response code other than 0 is a failure
}

$changes = array (
  array (
    'xpath'   => '/config/global/resources/default_setup/connection',
    'element' => 'host',
    'cdata'   => getenv( 'MYSQL_HOST' ) . ':' . getenv( 'MYSQL_PORT' ),
  ),
  array (
    'xpath'   => '/config/global/resources/default_setup/connection',
    'element' => 'dbname',
    'cdata'   => getenv( 'MYSQL_DATABASE' ),
  ),
  array (
    'xpath'   => '/config/global/resources/default_setup/connection',
    'element' => 'username',
    'cdata'   => getenv( 'MYSQL_USER' ),
  ),
  array (
    'xpath'   => '/config/global/resources/default_setup/connection',
    'element' => 'password',
    'cdata'   => getenv( 'MYSQL_PASSWORD' ),
  ),
  array (
    'xpath'   => '/config/global/resources/default_setup/connection',
    'element' => 'initStatements',
    'cdata'   => 'SET NAMES utf8',
  ),
  array (
    'xpath'   => '/config/global/resources/default_setup/connection',
    'element' => 'model',
    'cdata'   => 'mysql4',
  ),
  array (
    'xpath'   => '/config/global/resources/default_setup/connection',
    'element' => 'type',
    'cdata'   => 'pdo_mysql',
  ),
  array (
    'xpath'   => '/config/global/resources/default_setup/connection',
    'element' => 'pdoType',
    'cdata'   => '',
  ),
  array (
    'xpath'   => '/config/global',
    'element' => 'skip_process_modules_updates',
    'value'   => '1',
  ),
  array (
    'xpath'   => '/config/global/resources/db',
    'element' => 'table_prefix',
    'cdata'   => '',
  ),
  array (
    'xpath'   => '/config/global/crypt',
    'element' => 'key',
    'cdata'   => '',
    //@todo
  ),
  array (
    'xpath'   => '/config/global/install',
    'element' => 'date',
    'cdata'   => 'Mon, 25 Aug 2014 06:11:11 +0000',
  ),
  array (
    'xpath'   => '/config/admin/routers/adminhtml/args',
    'element' => 'frontName',
    'cdata'   => 'admin',
  ),
  array (
    'xpath'   => '/config/global',
    'element' => 'session_save',
    'cdata'   => 'files',
  ),
);

$trimmedPath   = rtrim($argv[$pathArg + 1], DIRECTORY_SEPARATOR);
$directoryPath = __DIR__ . $trimmedPath;
#$directoryPath = __DIR__ . DIRECTORY_SEPARATOR . 'app' . DIRECTORY_SEPARATOR . 'etc';
$templatePath  = $directoryPath . DIRECTORY_SEPARATOR . 'local.xml.template';
$finalPath     = $directoryPath . DIRECTORY_SEPARATOR . 'local.dev.xml';

try {
    $xml = file_get_contents( $templatePath );

    if (false === $xml) {
        fwrite(STDERR, "An error occurred." . PHP_EOL);
        exit(1); // A response code other than 0 is a failure
    }
} catch (Exception $e) {
    // Handle exception
    fwrite(STDERR, $e->getMessage() . PHP_EOL);
    exit(1); // A response code other than 0 is a failure
}

$dom = new DomDocument;
$dom->loadXML( $xml );
$xpath = new DOMXpath( $dom );

foreach ( $changes as $config ) {

    $connectionNodeList = $xpath->query( $config[ 'xpath' ] . '/' . $config[ 'element' ] );

    if ( isset( $config[ 'cdata' ] ) ) {
        $newNode = $dom->createElement( $config[ 'element' ] );
        $cData   = $dom->createCDATASection( $config[ 'cdata' ] );
        $newNode->appendChild( $cData );
    } else if ( isset( $config[ 'value' ] ) ) {
        $newNode = $dom->createElement( $config[ 'element' ], $config[ 'value' ] );
    }

    if ( $connectionNodeList->length == 0 ) {
        $connectionNodeList = $xpath->query( $config[ 'xpath' ] );
        $parentNode         = $connectionNodeList->item( 0 );
        $parentNode->appendChild( $newNode );
    } else {
        $oldNode = $connectionNodeList->item( 0 );
        $oldNode->parentNode->replaceChild( $newNode, $oldNode );
    }

}

$string = $dom->saveXML();

try {
    $success = file_put_contents( $finalPath, $string );

    if (false === $success) {
        fwrite(STDERR, "An error occurred." . PHP_EOL);
        exit(1); // A response code other than 0 is a failure
    }
} catch (Exception $e) {
    // Handle exception
    fwrite(STDERR, $e->getMessage() . PHP_EOL);
    exit(1); // A response code other than 0 is a failure
}
