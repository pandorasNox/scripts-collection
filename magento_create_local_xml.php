<?php

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

$directoryPath = __DIR__ . DIRECTORY_SEPARATOR . 'app' . DIRECTORY_SEPARATOR . 'etc';
$templatePath  = $directoryPath . DIRECTORY_SEPARATOR . 'local.xml.template';
$finalPath     = $directoryPath . DIRECTORY_SEPARATOR . 'local.dev.xml';

$xml = file_get_contents( $templatePath );

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

file_put_contents( $finalPath, $string );
