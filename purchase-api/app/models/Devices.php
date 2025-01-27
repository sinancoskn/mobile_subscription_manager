<?php

use Phalcon\Mvc\Model;
use Phalcon\Filter\Validation;
use Phalcon\Filter\Validation\Validator\PresenceOf;
use Phalcon\Filter\Validation\Validator\Numericality;
use Phalcon\Filter\Validation\Validator\InclusionIn;

class Devices extends Model
{

    /**
     *
     * @var integer
     */
    public $id;

    /**
     *
     * @var string
     */
    public $uid;

    /**
     *
     * @var integer
     */
    public $app_id;

    /**
     *
     * @var string
     */
    public $language;

    /**
     *
     * @var integer
     */
    public $os;

    /**
     *
     * @var string
     */
    public $created_at;

    /**
     * Initialize method for model.
     */
    public function initialize()
    {
        $this->setSchema("public");
        $this->setSource("devices");
    }

    /**
     * Allows to query a set of records that match the specified conditions
     *
     * @param mixed $parameters
     * @return Devices[]|Devices|\Phalcon\Mvc\Model\ResultSetInterface
     */
    public static function find($parameters = null): \Phalcon\Mvc\Model\ResultsetInterface
    {
        return parent::find($parameters);
    }

    /**
     * Allows to query the first record that match the specified conditions
     *
     * @param mixed $parameters
     * @return Devices|\Phalcon\Mvc\Model\ResultInterface|\Phalcon\Mvc\ModelInterface|null
     */
    public static function findFirst($parameters = null): ?\Phalcon\Mvc\ModelInterface
    {
        return parent::findFirst($parameters);
    }

    public function beforeCreate()
    {
        $this->created_at = date('Y-m-d H:i:s');
    }

    public function validation()
    {
        $validator = new Validation();

        // Validate app_id (must be a number)
        $validator->add(
            'app_id',
            new Numericality([
                'message' => 'The app_id must be a number.',
            ])
        );

        // Validate language (required)
        $validator->add(
            'language',
            new PresenceOf([
                'message' => 'The language is required.',
            ])
        );

        // Validate OS (must be 1 or 2)
        $validator->add(
            'os',
            new InclusionIn([
                'domain'  => [1, 2],
                'message' => 'The os must be either 1 (iOS) or 2 (Android).',
            ])
        );

        if (empty($this->created_at)) {
            $this->created_at = date('Y-m-d H:i:s');
        }

        return $this->validate($validator);
    }

}
